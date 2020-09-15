import router from '../router';
import {message} from 'ant-design-vue';
import moment from "moment";

export function trimAll(str) {
    let result;
    result = str.replace(/(^\s+)|(\s+$)/g, "");
    result = result.replace(/\s/g, "");
    return result;
}

export function errors(error) {
    console.log(error)
    if ((error + '').indexOf('Network Error') != -1) {
        message.warning("网络错误请重试!");
        return false;
    } else if ((error + '').indexOf('timeout') != -1) {
        message.warning("网络超时请重试!");
        return false;
    }
    switch (error.response.status) {
        case 400: {
            message.warning(error.response.data.error_description);
            break;
        }
        case 401: {
            message.warning("身份验证失败，请重新登陆");
            router.push({
                path: '/login',
                query: {
                    // redirect: router.currentRoute.fullPath
                }  //从哪个页面跳转
            });
            break;
        }
        case 404: {
            message.warning("哎呀，404了！杀个程序员祭天吧！");
            break;
        }
        case 500: {
            message.warning("哎呀，出错了！杀个程序员祭天吧！");
            break;
        }
    }
}

//返回的是对象形式的参数
export function getUrlArgObject(parm1) {
    let args = new Object();
    let query = window.location.href;//获取查询串
    let pairs;
    if (query.indexOf("?") != -1) {
        pairs = query.split("?")[1].split("&");
        for (let i = 0; i < pairs.length; i++) {
            let pos = pairs[i].indexOf('=');//查找name=value
            if (pos == -1) {//如果没有找到就跳过
                continue;
            }
            let argname = pairs[i].substring(0, pos);//提取name
            let value = decodeURIComponent(pairs[i].substring(pos + 1));//提取value
            args[argname] = unescape(value);//存为属性
        }
    }
    return args[parm1];//返回对象
}

export function momentDate(t) {
    if (t === undefined) {
        return null;
    }

    const date = new Date(Number(t) * 1000);
    const d = moment(date, 'YYYY-MM-DD');

    return d;
}

export function DateTimeFormat(t, format) {
    if (t === undefined) {
        return null;
    }

    const date = new Date(Number(t) * 1000);
    if (format === undefined || format === '') {
        format = "YYYY-MM-DD HH:mm:ss"
    }
    return moment(date).format(format);
}