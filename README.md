# showks-github-repository-operator

## Install

```
$ echo "GITHUB_TOKEN=xxxxxxxxxxz" > ./config/secret.env
```

```
$ make deploy
```

##sage

```yaml
apiVersion: showks.cloudnativedays.jp/v1beta1
kind: GithubRepository
metadata:
  name: sample
spec:
  org: cloudnativedaysjp
  name: showks-canvas-sampleusername
  repositoryTemplate:
    org: containerdaysjp
    name: showks-canvas
    #　テンプレートリポジトリからチェックアウトして作成するブランチを指定します
    initialBranches:
      - refs/heads/master:refs/heads/master
      - refs/heads/master:refs/heads/staging
      - refs/heads/master:refs/heads/feature
  collaborators:
    - name: alice
      permission: admin
    - name: bob
      permission: pull
  branchProtection:
    enforceAdmins: false
    requiredPullRequestReviews: nil
    requiredStatusChecks:
      strict: true
      contexts: []
    restrictions:
      users: []
      teams:
        - showks-members
  webhooks:
    - name: web
      config:
        url: https://example.com
        contentType: json 
      events:
        - push
      active: true
```

## Development

GithubのAPIを実行するため、環境変数`GITHUB_TOKEN`にパーソナルアクセストークンをセットします。Scopeには`repo`と`user`、`delete_repo`を設定します。

```
$ export GITHUB_TOKEN=xxxxxxxxxx
```

以下のようにコントローラーを手元で実行します。

```
$ kubectl apply -f config/crds
$ make run
```

