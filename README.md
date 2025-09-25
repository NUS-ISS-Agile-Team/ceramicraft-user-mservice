# ceramicraft-user-mservice
## Functional Features
1. support customer/platform admin login/logout, same API with cookie setted under different domain
2. authenticate with jwt token
3. support new customer registration, verified by email opt code
4. provide auth filter for other microservice use

## Detail Disign
[see the document](https://cerami-craft.atlassian.net/wiki/spaces/swe5006gro/pages/4554754/UserService)

## startup dependency
1. Mysql
2. email sending config
3. docker startup
```
cd server
docker build -t ${DOCKER_HUB_USERNAME}/ceramicraft-user-mservice:<version> .
docker run -d --name ceramicraft-user-mservice -e MYSQL_PASSWORD=${{ secrets.MYSQL_PASSWORD }} -e SMTP_PASSWORD=${{ secrets.SMTP_PASSWORD }} -e SMTP_EMAIL_FROM=${{ secrets.SMTP_EMAIL_FROM }} "${DOCKER_HUB_USERNAME}/ceramicraft-user-mservice:${{ github.event.inputs.version }}"
```
4. non-containerized startup
```
MYSQL_PASSWORD=your_password SMTP_PASSWORD=your_smtp_password SMTP_EMAIL_FROM=your_email@example.com go run server/main.go
```