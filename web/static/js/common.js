/**
 * CopilotLens 公共前端模块
 * 提供 SSR 初始数据读取、sessionStorage 缓存、后台静默刷新、loading 控制。
 */
(function (global) {
    'use strict';

    var CL = {
        _initialDataParsed: false
    };

    function parseInitialData() {
        if (CL._initialDataParsed) return;
        CL._initialDataParsed = true;
        global.__INITIAL_DATA__ = {};
        var el = document.getElementById('initial-data');
        if (!el || !el.textContent) return;
        try {
            global.__INITIAL_DATA__ = JSON.parse(el.textContent) || {};
        } catch (e) {
            console.warn('[CL] 解析 initial-data 失败:', e);
            global.__INITIAL_DATA__ = {};
        }
    }

    CL.getInitialData = function (key) {
        parseInitialData();
        if (!key) return global.__INITIAL_DATA__ || {};
        var data = global.__INITIAL_DATA__ || {};
        return data.hasOwnProperty(key) ? data[key] : null;
    };

    CL.setLoading = function (show) {
        var el = document.getElementById('loadingBox');
        if (!el) return;
        el.style.display = show ? 'flex' : 'none';
    };

    CL.showNoData = function (show) {
        var el = document.getElementById('noData');
        if (!el) return;
        el.style.display = show ? 'block' : 'none';
    };

    CL.formatNumber = function (n) {
        var num = Number(n);
        if (!isFinite(num)) num = 0;
        return num.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
    };

    function storageKeyFor(url, params) {
        var qs = [];
        if (params) {
            var keys = Object.keys(params).sort();
            for (var i = 0; i < keys.length; i++) {
                qs.push(keys[i] + '=' + encodeURIComponent(params[keys[i]]));
            }
        }
        return url + '?' + qs.join('&');
    }

    function readSession(key) {
        try {
            var raw = sessionStorage.getItem(key);
            if (!raw) return null;
            var obj = JSON.parse(raw);
            if (obj && obj._ts && (Date.now() - obj._ts > 30 * 60 * 1000)) {
                sessionStorage.removeItem(key);
                return null;
            }
            return obj ? obj.data : null;
        } catch (e) {
            return null;
        }
    }

    function writeSession(key, data) {
        try {
            sessionStorage.setItem(key, JSON.stringify({ _ts: Date.now(), data: data }));
        } catch (e) {
            console.warn('[CL] sessionStorage 写入失败:', e);
        }
    }

    /**
     * 本地时间格式化辅助（避免 UTC 偏差）
     */
    CL.localDateStr = function (d) {
        var year = d.getFullYear();
        var month = String(d.getMonth() + 1).padStart(2, '0');
        var day = String(d.getDate()).padStart(2, '0');
        return year + '-' + month + '-' + day;
    };
    CL.localMonthStr = function (d) {
        var year = d.getFullYear();
        var month = String(d.getMonth() + 1).padStart(2, '0');
        return year + '-' + month;
    };

    CL.toggleLang = function () {
        var lang = localStorage.getItem('lang') === 'en' ? 'zh' : 'en';
        localStorage.setItem('lang', lang);
        location.reload();
    };

    CL.applyLang = function () {
        $('[data-i18n]').each(function () {
            var key = $(this).data('i18n');
            var val = key.split('.').reduce(function (o, k) { return o && o[k]; }, window.i18n);
            if (val) $(this).text(val);
        });
    };

    CL.updateLangSwitch = function () {
        var isEn = localStorage.getItem('lang') === 'en';
        $('.lang-zh').toggleClass('lang-active', !isEn).toggleClass('lang-inactive', isEn);
        $('.lang-en').toggleClass('lang-active', isEn).toggleClass('lang-inactive', !isEn);
    };

    /**
     * 带 sessionStorage 的 stale-while-revalidate 数据请求。
     * @param {string} url
     * @param {object} params
     * @param {object} options
     *   - onSuccess(data): 拿到数据时回调（包括缓存命中时）
     *   - onError(err): 请求失败时回调
     *   - onRender?(data): 首次渲染回调（可选，默认调用 onSuccess）
     *   - storageKey?: 自定义 sessionStorage key
     */
    CL.fetchWithCache = function (url, params, options) {
        options = options || {};
        var key = options.storageKey || storageKeyFor(url, params);
        var cached = readSession(key);
        var rendered = false;

        function render(data) {
            rendered = true;
            if (typeof options.onRender === 'function') {
                options.onRender(data);
            } else if (typeof options.onSuccess === 'function') {
                options.onSuccess(data);
            }
        }

        if (cached) {
            CL.showNoData(false);
            render(cached);
        }

        $.get(url, params)
            .done(function (data) {
                writeSession(key, data);
                if (!rendered) {
                    CL.setLoading(false);
                    CL.showNoData(false);
                    render(data);
                } else if (typeof options.onSuccess === 'function') {
                    options.onSuccess(data);
                }
            })
            .fail(function (err) {
                if (!rendered) {
                    CL.setLoading(false);
                    CL.showNoData(true);
                }
                if (typeof options.onError === 'function') {
                    options.onError(err);
                }
            });
    };

    global.CL = CL;
})(window);
