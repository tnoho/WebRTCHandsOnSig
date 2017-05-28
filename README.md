# SwiftでWebRTC実装ハンズオン用 Signaling server

このSignalingサーバは2017/06/18に開催される[集まれSwift好き！Swift愛好会 vol20 \~SwiftでWebRTC実装ハンズオン\~](https://love-swift.connpass.com/event/55249/)の為に作成しました。
 
ハンズオン向けにメッセージ内で工夫しなくとも、他人のSignalingと混信しにくいようにURLの末尾でWebSocketのRoomが変わるようになっています。
 
はじめてgoで書きました。おかしな点があればご指摘ください。

# 謝辞

templateに利用したのは [yusuke84](https://github.com/yusuke84) 氏の [WebRTCハンズオン](http://qiita.com/yusuke84/items/de9f0f6d221acec6fc07) 用ソースコードです。ありがとうございます！
 
 作成に当たり [matryer](https://github.com/matryer) 氏の [goblueprint](https://github.com/matryer/goblueprints) の和訳である、[Go言語によるWebアプリケーション開発](https://www.oreilly.co.jp/books/9784873117522/) のチュートリアルをもとに書きました。筆者と訳者に感謝します。

 # デプロイ方法

 ## ubuntuの場合

signaling.service.template の中身を自分の展開先のPATHに変更して`/lib/systemd/system/signaling.service`にmv

 > sudo mv signaling.service /lib/systemd/system/ 

 下記のコマンドで登録して

 > sudo systemctl enable signaling.service 

 下記のコマンドで実行

 > sudo systemctl start signaling 

 ### 注意

 getUserMediaはSSL必須のため、nginxのssl reverse proxy化で動かす事を想定しています。reverse proxy時にHostの書換をお忘れ無く。