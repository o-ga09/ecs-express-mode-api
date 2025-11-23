# AWS OIDC連携によるECRへのデプロイ手順

このドキュメントでは、GitHub ActionsからOIDC認証を使用してAWS ECRにDockerイメージをプッシュする手順を説明します。

## 概要

- GitHub ActionsとAWS間でOIDC連携を使用することで、長期的なアクセスキーを使わずに安全にAWSリソースにアクセスできます
- ビルドしたDockerイメージをAmazon ECRにプッシュします

## 前提条件

- AWSアカウント
- GitHubリポジトリの管理者権限
- AWS CLIがインストールされていること（設定用）

## 手順

### 1. AWS側の設定

#### 1.1 OIDCプロバイダーの作成

AWSコンソールまたはCLIでGitHub ActionsのOIDCプロバイダーを作成します。

**AWS CLIの場合:**

```bash
aws iam create-open-id-connect-provider \
  --url https://token.actions.githubusercontent.com \
  --client-id-list sts.amazonaws.com \
  --thumbprint-list 6938fd4d98bab03faadb97b34396831e3780aea1
```

**コンソールの場合:**

1. IAMコンソール → IDプロバイダー → プロバイダーを追加
2. プロバイダーのタイプ: `OpenID Connect`
3. プロバイダーのURL: `https://token.actions.githubusercontent.com`
4. 対象者: `sts.amazonaws.com`

#### 1.2 IAMロールの作成

GitHub ActionsがAssumeRoleできるIAMロールを作成します。

**信頼ポリシー:**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<AWS_ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:<GITHUB_USERNAME>/<REPO_NAME>:ref:refs/heads/main"
        }
      }
    }
  ]
}
```

**置き換えが必要な項目:**
- `<AWS_ACCOUNT_ID>`: あなたのAWSアカウントID
- `<GITHUB_USERNAME>`: GitHubのユーザー名または組織名
- `<REPO_NAME>`: リポジトリ名（例: `ecs-express-mode-api`）

**AWS CLIでロールを作成:**

```bash
# 信頼ポリシーをファイルに保存
cat > trust-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::<AWS_ACCOUNT_ID>:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com"
        },
        "StringLike": {
          "token.actions.githubusercontent.com:sub": "repo:<GITHUB_USERNAME>/<REPO_NAME>:ref:refs/heads/main"
        }
      }
    }
  ]
}
EOF

# ロールを作成
aws iam create-role \
  --role-name GitHubActionsECRRole \
  --assume-role-policy-document file://trust-policy.json
```

#### 1.3 ECRへのアクセス権限を付与

作成したロールにECRへのアクセス権限を付与します。

```bash
aws iam attach-role-policy \
  --role-name GitHubActionsECRRole \
  --policy-arn arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser
```

**カスタムポリシーを使う場合:**

より細かい権限制御が必要な場合は、以下のようなカスタムポリシーを作成します。

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:PutImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload"
      ],
      "Resource": "*"
    }
  ]
}
```

#### 1.4 ECRリポジトリの作成

Dockerイメージを保存するECRリポジトリを作成します。

```bash
aws ecr create-repository \
  --repository-name ecs-express-mode-api \
  --region ap-northeast-1
```

### 2. GitHub側の設定

#### 2.1 GitHub Secretsの設定

リポジトリのSettings → Secrets and variables → Actions → New repository secretで以下を追加:

- **Secret名**: `AWS_ROLE_ARN`
- **値**: 作成したIAMロールのARN（例: `arn:aws:iam::123456789012:role/GitHubActionsECRRole`）

ロールのARNは以下のコマンドで確認できます:

```bash
aws iam get-role --role-name GitHubActionsECRRole --query 'Role.Arn' --output text
```

### 3. ワークフローファイルの確認

`.github/workflows/deploy.yml`が以下の内容になっていることを確認します:

```yaml
name: Build and Push to ECR

on:
  push:
    branches:
      - main
  workflow_dispatch:

env:
  AWS_REGION: ap-northeast-1
  ECR_REPOSITORY: ecs-express-mode-api

permissions:
  id-token: write
  contents: read

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE_ARN }}
          aws-region: ${{ env.AWS_REGION }}

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build, tag, and push image to Amazon ECR
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          docker buildx build \
            --platform linux/amd64 \
            --tag $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG \
            --tag $ECR_REGISTRY/$ECR_REPOSITORY:latest \
            --push \
            .

      - name: Output image URI
        env:
          ECR_REGISTRY: ${{ steps.login-ecr.outputs.registry }}
          IMAGE_TAG: ${{ github.sha }}
        run: |
          echo "Image URI: $ECR_REGISTRY/$ECR_REPOSITORY:$IMAGE_TAG"
```

### 4. デプロイの実行

#### 4.1 自動デプロイ

`main`ブランチにプッシュすると自動的にワークフローが実行されます。

```bash
git add .
git commit -m "Add GitHub Actions workflow for ECR push"
git push origin main
```

#### 4.2 手動デプロイ

GitHubリポジトリのActions タブから手動で実行することもできます。

1. Actionsタブを開く
2. "Build and Push to ECR" ワークフローを選択
3. "Run workflow" ボタンをクリック

### 5. 動作確認

#### 5.1 GitHub Actionsのログ確認

GitHub上でワークフローの実行ログを確認し、エラーがないことを確認します。

#### 5.2 ECRリポジトリの確認

ECRにイメージがプッシュされたことを確認します。

```bash
aws ecr describe-images \
  --repository-name ecs-express-mode-api \
  --region ap-northeast-1
```

または、AWSコンソールのECRページで確認できます。

## トラブルシューティング

### エラー: "User is not authorized to perform: sts:AssumeRoleWithWebIdentity"

- IAMロールの信頼ポリシーが正しく設定されているか確認
- GitHubリポジトリ名が信頼ポリシーと一致しているか確認
- OIDCプロバイダーが正しく作成されているか確認

### エラー: "No basic auth credentials"

- `AWS_ROLE_ARN`がGitHub Secretsに正しく設定されているか確認
- ロールにECRへのアクセス権限が付与されているか確認

### エラー: "Repository does not exist"

- ECRリポジトリが作成されているか確認
- リポジトリ名が`ECR_REPOSITORY`環境変数と一致しているか確認
- リージョンが正しいか確認

## カスタマイズ

### リージョンの変更

`.github/workflows/deploy.yml`の`AWS_REGION`環境変数を変更します。

```yaml
env:
  AWS_REGION: us-east-1  # 任意のリージョンに変更
  ECR_REPOSITORY: ecs-express-mode-api
```

### リポジトリ名の変更

`.github/workflows/deploy.yml`の`ECR_REPOSITORY`環境変数を変更します。

```yaml
env:
  AWS_REGION: ap-northeast-1
  ECR_REPOSITORY: your-repository-name  # 任意のリポジトリ名に変更
```

### 特定のブランチでのみ実行

`.github/workflows/deploy.yml`の`on.push.branches`を変更します。

```yaml
on:
  push:
    branches:
      - main
      - develop
      - production
```

## セキュリティのベストプラクティス

1. **最小権限の原則**: IAMロールには必要最小限の権限のみを付与
2. **ブランチ制限**: 信頼ポリシーで特定のブランチからのみAssumeRoleできるように制限
3. **シークレットの管理**: AWS_ROLE_ARNはGitHub Secretsで管理し、コードにハードコードしない
4. **定期的な監査**: CloudTrailでAssumeRoleの使用状況を監視

## 参考リンク

- [GitHub ActionsでのOIDC連携 - AWS公式ドキュメント](https://docs.github.com/ja/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-amazon-web-services)
- [Amazon ECR - AWS公式ドキュメント](https://docs.aws.amazon.com/ja_jp/ecr/)
- [aws-actions/configure-aws-credentials](https://github.com/aws-actions/configure-aws-credentials)
- [aws-actions/amazon-ecr-login](https://github.com/aws-actions/amazon-ecr-login)
