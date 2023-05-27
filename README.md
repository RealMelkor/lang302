# lang302

FCGI programs to redirect users to the most appropriate pages based on their
HTTP header Accept-Language value.

Configuration example for [nginx](https://nginx.org/) :

```nginx
location ~ ^/$  {
    fastcgi_pass    127.0.0.1:9000;
    include         fastcgi_params;
}
```

Configuration example for [OpenBSD httpd](https://man.openbsd.org/httpd.8) :

```nginx
location "/" {
    fastcgi socket tcp localhost 9000
}
```

With the default configuration file, this will redirect requests on / to /en or
/fr depending on the value of Accept-Language.
