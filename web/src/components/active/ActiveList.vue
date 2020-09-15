<template>
    <div>
        <h3 align="left">活动 > 列表</h3>
        <a-form-model layout="inline" align="left" :style="{width:'100%'}">
            <a-form-model-item>
                <a-radio-group v-model="prefix" @change="radioChange">
                    <a-radio value="master">
                        管理节点
                    </a-radio>
                    <a-radio value="worker">
                        工作节点
                    </a-radio>
                    <a-radio value="flow">
                        工作流
                    </a-radio>
                    <a-radio value="job">
                        任务
                    </a-radio>
                </a-radio-group>
            </a-form-model-item>
            <a-form-model-item>
                <a-input placeholder="Key" v-model="key"></a-input>
            </a-form-model-item>
            <a-form-model-item>
                <a-button type="primary" @click="filter">筛选</a-button>
            </a-form-model-item>
        </a-form-model>

        <a-table
                :columns="columns"
                :rowKey="record => record.path"
                :dataSource="data"
                :loading="loading"
                :pagination="false"
        >

            <span slot="action">

            </span>
        </a-table>

    </div>
</template>

<script>
    import axios from "../../library/http";
    import {DateTimeFormat} from "../../library/utils";

    export default {
        name: "ActiveList",
        mounted() {
            this.fetch();
        },
        data() {
            return {
                data: [],
                loading: false,
                columns: [
                    {
                        title: 'Key',
                        dataIndex: 'path',
                    },
                    {
                        title: 'Value',
                        dataIndex: 'value',
                    },
                    {
                        title: '操作',
                        key: 'action',
                        fixed: 'right',
                        width: 180,
                        scopedSlots: {
                            customRender: 'action'
                        },
                    },
                ],
                prefix: 'master',
                key: ''
            };
        },
        methods: {
            fetch() {
                if (this.$route.query.prefix !== undefined && this.$route.query.prefix !== '') {
                    this.prefix = this.$route.query.prefix;
                }
                if (this.$route.query.key !== undefined && this.$route.query.key !== '') {
                    this.key = this.$route.query.key;
                }
                this.loading = true;
                axios({
                    url: '/etcd/list',
                    data: {
                        client: 'crontab',
                        prefix: this.prefix,
                        key: this.key,
                    },
                    type: 'json',
                }).then(data => {
                    console.log(data)
                    this.loading = false;
                    this.data = data.data.keys;
                });
            },
            detail(id) {
                this.$router.push("/job/detail?id=" + id);
            },
            filter() {
                let para = [];
                if (this.prefix !== undefined && this.prefix !== '') {
                    para = [...para, 'prefix=' + this.prefix]
                }
                if (this.key !== undefined && this.key !== '') {
                    para = [...para, 'key=' + this.key]
                }
                let url = "/active/list";
                if (para.length > 0) {
                    url = url + "?" + para.join('&');
                }
                if (this.$route.fullPath !== url) {
                    this.$router.push(url);
                } else {
                    location.reload();
                }
            },
            DateTimeFormat,
            radioChange(e) {
                console.log('radio checked', e.target.value);
                this.prefix = e.target.value;
            },
        },
    }

</script>

<style scoped>
    .highlight {
        background-color: rgb(255, 192, 105);
        padding: 0px;
    }
</style>