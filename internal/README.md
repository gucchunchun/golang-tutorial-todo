# ```internal```

非公開アプリケーションとライブラリーのコード部分。他人にimportされたくないコードをこの中に記載する。<br/>
どの階層内においても```/internal```ディレクトリは同じ用法で使用されることがある。

- アプリケーション自体のコード：```/internal/app```
- アプリケーション内で共通で使用されているコード：```/internal/pkg```
- 同じ親ディレクトリ内にない限り、importすることはできない

[Reference](https://github.com/golang-standards/project-layout?tab=readme-ov-file)
