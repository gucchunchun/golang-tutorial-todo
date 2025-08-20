# ```cmd```

プロジェクトのメイン部分。<br/>
```/internal```または```/pkg```にあるコードを使用するのに主に使用される。

- このファイル内にたくさんのコードを書くべきではない
- 大抵```main()```を持つ

※再利用可能な内容は```/pkg```、他人に利用されたくない内容は```/internal```に移動するべき

[Reference](https://github.com/golang-standards/project-layout/blob/master/cmd/README.md)
