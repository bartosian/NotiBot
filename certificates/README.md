# Configuring Nginx Proxy

These instructions will guide you through the process of configuring an Nginx reverse proxy on Ubuntu.

## Prerequisites

- Ubuntu server with sudo access
- Basic understanding of Nginx configuration files

## Installation

1. Update the package index on your server:

```shell
sudo apt update
```

2. Install Nginx:

```shell
sudo apt install nginx
```

3. Verify that Nginx is running by checking its status:

```shell
systemctl status nginx
```

## Configuration

4. Create a new Nginx configuration file for your domain:

```shell
sudo vim /etc/nginx/sites-available/your_domain
```

5. Use the example configuration in the `nginx-config-example` file provided, and modify it to suit your needs.

6. Create a symbolic link to enable your new configuration:

```shell
sudo ln -s /etc/nginx/sites-available/your_domain /etc/nginx/sites-enabled/
```

7. Verify that your configuration file is valid:

```shell
sudo nginx -t
```

8. If there are no errors, restart Nginx to apply your new configuration:

```shell
sudo systemctl restart nginx
```

9. Install Certbot and the Nginx plugin:

```shell
sudo apt install certbot python3-certbot-nginx
```

10. Obtain SSL/TLS certificates and configure Nginx to use them by running the following command and following the prompts:

```shell
sudo certbot --nginx -d example.com -d www.example.com
```

Note: Replace example.com with your actual domain name.

When prompted, choose option 2 to redirect traffic.

11. Verify that Certbot's automatic renewal service is active and running:

```shell
sudo systemctl status certbot.timer
```

12. Test the renewal process by running the following command (the --dry-run option tells Certbot to simulate the renewal process without actually renewing the certificates):

```shell
sudo certbot renew --dry-run
```

For more information, check out the following tutorials on DigitalOcean:

- [How To Configure Nginx as a Reverse Proxy on Ubuntu 22.04](https://www.digitalocean.com/community/tutorials/how-to-configure-nginx-as-a-reverse-proxy-on-ubuntu-22-04)
- [How To Secure Nginx with Let's Encrypt on Ubuntu 20.04](https://www.digitalocean.com/community/tutorials/how-to-secure-nginx-with-let-s-encrypt-on-ubuntu-20-04)
