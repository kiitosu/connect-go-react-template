# Product Context

## なぜこのプロジェクトが存在するのか
- GoバックエンドとReactフロントエンド間でgRPC通信を行う際、セットアップやAPIスキーマ管理、開発環境構築に多くの手間がかかる。
- connect-go/connect-webの導入・運用ノウハウが分散しており、ベストプラクティスが見えづらい。
- プロトコルバッファやDockerを活用したモダンなフルスタック開発の標準構成が求められている。

## 解決する課題
- gRPC+connectプロトコルによるAPI連携の導入障壁を下げる
- バックエンド・フロントエンドの分離開発を容易にする
- APIスキーマの一元管理と自動生成の仕組みを提供
- Docker Composeによる環境構築の簡易化
- 新規プロジェクト立ち上げ時の初期構成検討の手間を削減

## 理想的なユーザー体験
- 最小限のセットアップでgRPC通信が動作する
- APIスキーマの変更が即座に各層へ反映される
- バックエンド・フロントエンドの開発が独立して進められる
- Docker Composeで一発起動・一発停止が可能
- サンプルAPIを通じて実装イメージを素早く掴める

## 想定ユーザー
- gRPC/Protocol Buffers/Connectの導入を検討しているフルスタックエンジニア
- モダンなAPI連携のベストプラクティスを学びたい開発者
- 新規プロジェクトの雛形を探しているチーム
