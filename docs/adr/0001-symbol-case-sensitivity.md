# ADR 0001: シンボルの大文字小文字の扱い

## ステータス

採用（Accepted）

## 日付

2025-10-13

## コンテキスト

Common Lispでは、シンボルはデフォルトで大文字に変換される。つまり`foo`と`FOO`は同じシンボルとして扱われる。一方、GoLispはGo言語との相互運用を目的としており、Go関数名やstructフィールド名は大文字小文字を区別する。

この設計判断では、GoLispがシンボルの大文字小文字をどう扱うべきかを決定する。

## 検討した選択肢

### オプション1: 大文字に統一（Common Lisp方式）

シンボルを読み込み時に大文字に変換する。

**メリット:**
- Common Lisp互換性が高い
- 既存のCommon Lispコードがそのまま動く
- Lisp文化に沿っている
- 多くのLisp書籍・チュートリアルと一致

**デメリット:**
- Go連携が煩雑になる
  - `fmt.Println` → `FMT.PRINTLN`（不自然）
  - Go structフィールド: `UserName` → `USERNAME`（不一致）
- 現代的な言語トレンドと乖離
- JSON/Web APIとの統合が困難
  - `{"userName": "Alice"}` → `(gethash 'USERNAME json)` （キー名不一致）
- 外部システム統合時に常に変換が必要

### オプション2: 小文字を区別（モダン方式）

シンボルの大文字小文字をそのまま保持し、区別する。

**メリット:**
- Go連携がスムーズ
  - `fmt.Println` をそのまま呼べる
  - structフィールド名が一致: `(go:field user 'UserName)`
- モダンな言語（JavaScript, Python, Rust, Go）と同じ動作
- 外部システム統合が自然
  - JSONキー名、HTTPヘッダー名、DBカラム名がそのまま使える
- 型システムとの相性が良い
  - Goのstructとのマッピングが直接的

**デメリット:**
- Common Lisp互換性が下がる
- 既存のCommon Lispコードを移植する際に変更が必要
- `defun` と `DEFUN` は別のシンボルとして扱われる

## 決定

**オプション2（小文字を区別）を採用する。**

## 理由

1. **プロジェクトの目的との整合性**
   - 仕様書（SPECIFICATION.md:6）で明記されている「Go言語との相互運用」が最優先
   - M9（Go連携）、M11（JSON/HTTP）での実用性を重視

2. **モダンな開発環境との親和性**
   - 現代のプログラミング言語はすべて大文字小文字を区別
   - Web API、JSON、データベースなど外部システムとの統合がスムーズ

3. **実用性**
   - Go関数を自然に呼べる: `(fmt.Println "hello")`
   - structフィールドアクセスが直接的: `(go:field user 'UserName)`
   - JSONパースが自然: `(gethash 'userName json)`

4. **将来の拡張性**
   - Go連携機能（M9）での実装が簡潔
   - HTTP/JSON機能（M11以降）での利便性が高い

## 影響

### 肯定的影響

- Go関数呼び出しがそのまま書ける
- 外部システムとのデータ交換が簡潔
- モダンな開発スタイルとの整合性

### 否定的影響

- Common Lispコードを直接実行できない
- Common Lisp書籍のコード例をそのまま写経できない場合がある

### 移行戦略

Common Lisp互換性が必要な場合の対応策:

1. **ヘルパー関数の提供**
```lisp
(defun symbol-equal-ignore-case (s1 s2)
  (string-equal (symbol-name s1) (symbol-name s2)))
```

2. **リーダーマクロのオプション機能（将来的に）**
```lisp
;; オプション機能として実装可能
(with-cl-reader-case :upcase
  (defun square (x) (* x x)))  ; SQUAREとして定義される
```

3. **ドキュメントでの明示**
   - GoLispはシンボルの大小文字を区別することを明記
   - Common Lispコード移植時のガイドラインを提供

## 参照

- SPECIFICATION.md:28-30（シンボルの大小文字の扱い）
- M9マイルストーン（Go連携機能）
- M11マイルストーン（JSON/HTTP機能）

## 備考

この決定は、Common Lisp純粋主義よりも実用性を優先するGoLispの設計方針を反映している。Common Lispの思想を学びつつ、Go言語エコシステムとの統合を重視するという、プロジェクトのバランスを取った判断である。
