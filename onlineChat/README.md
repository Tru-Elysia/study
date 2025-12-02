# 使用教程

## 启动服务器

* **双击server.exe以启动服务器**

## 登录

* **POST ： <http://localhost:8080/login>**
* **body ：**

    |   参数名    |   类型   |  参数值  |
    | :------: | :----: | :---: |
    | username | string | （用户名） |

## 登出

* **POST ： <http://localhost:8080/logout>**
* **body ：**

    |   参数名    |   类型   |  参数值  |
    | :------: | :----: | :---: |
    | username | string | （用户名） |

## 发送消息

* **POST ： <http://localhost:8080/send>**
* **body ：**

    |   参数名    |   类型   |  参数值  |
    | :------: | :----: | :---: |
    | username | string | （用户名） |
    | message  | string | （消息）  |

## 聊天记录

* **GET ： <http://localhost:8080//message>**

## 在线用户

* **GET ： <http://localhost:8080//onlineuser>**

