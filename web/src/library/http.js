import axios from 'axios'
import Cookies from 'js-cookie'
import {errors} from "./utils";
import {message} from 'ant-design-vue';
import JSONBIG from 'json-bigint'

// axios 配置
axios.defaults.timeout = 60000;
axios.defaults.method = 'POST';
axios.defaults.headers.common['Content-Type'] = 'application/json';
axios.defaults.baseURL = process.env.VUE_APP_API_URL;
axios.defaults.transformResponse = [data => {
    try {
        data = JSONBIG.parse(data);
    } catch (e) {
        //console.log(e)
    }
    return data;
}];

// http request 拦截器，通过这个，我们就可以把Cookie传到后台
axios.interceptors.request.use(
    config => {
        const token = Cookies.get('token'); //获取Cookie
        if (token !== undefined) {
            config.headers['Authorization'] = 'Bearer ' + token; //后台接收的参数，后面我们将说明后台如何接收
        }
        return config;
    },
    err => {
        errors(err);
        return Promise.reject(err);
    }
);


// http response 拦截器
axios.interceptors.response.use(
    response => {
        if ((response.data) && response.data.code && response.data.code != 10000) {
            message.warning(response.data.msg);
        }
        return response.data;
    },
    error => {
        errors(error, this);
        return Promise.reject(error)
    });

export default axios;