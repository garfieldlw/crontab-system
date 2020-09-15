<template>
    <div>
        <h3 align="left">任务 > 列表</h3>
        <a-form-model layout="inline" align="left" :style="{width:'100%'}">
            <a-form-model-item>
                <a-input placeholder="Id" v-model="id"></a-input>
            </a-form-model-item>
            <a-form-model-item>
                <a-input placeholder="名称" v-model="name"></a-input>
            </a-form-model-item>
            <a-form-model-item>
                <a-button type="primary" @click="filter">筛选</a-button>
            </a-form-model-item>
        </a-form-model>

        <a-table
                :columns="columns"
                :rowKey="record => record.id.toString()"
                :dataSource="data"
                :pagination="pagination"
                :loading="loading"
                @change="handleTableChange"
                :scroll="{ x:  'calc(600px + 40%)' }"
        >

            <span slot="time" slot-scope="text, record">{{DateTimeFormat(record.create_time)}}</span>

            <span slot="action" slot-scope="text, record">
                <a @click="() => detail(record.id.toString())">详情</a>
            </span>
        </a-table>

    </div>
</template>

<script>
    import axios from "../../library/http";
    import {DateTimeFormat} from "../../library/utils";

    export default {
        name: "JobList",
        mounted() {
            this.fetch();
        },
        data() {
            return {
                dataClient: [],
                data: [],
                pagination: {},
                loading: false,
                columns: [
                    {
                        title: 'Id',
                        dataIndex: 'id',
                        width: 200,
                    },
                    {
                        title: '名称',
                        dataIndex: 'name',
                        width: 250,
                    },
                    {
                        title: '类型',
                        dataIndex: 'job_type',
                    },
                    {
                        title: '1启用，2禁用',
                        dataIndex: 'status',
                        width: 200,
                    },
                    {
                        title: '创建时间',
                        dataIndex: 'create_time',
                        width: 200,
                        scopedSlots: {
                            customRender: 'time'
                        },
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
                id: '',
                name: '',
            };
        },
        methods: {
            handleTableChange(pagination) {
                console.log(pagination);
                const pager = {...this.pagination};
                pager.current = pagination.current;
                this.pagination = pager;
                this.fetch({
                    page: pagination.current,
                });
            },
            fetch(params = {}) {
                if (this.$route.query.id !== undefined && this.$route.query.id !== '') {
                    this.id = this.$route.query.id;
                }
                if (this.$route.query.name !== undefined && this.$route.query.name !== '') {
                    this.name = this.$route.query.name;
                }
                console.log('params:', params);
                this.loading = true;
                let offset = 0;
                if (params.page > 0) {
                    offset = (params.page - 1) * 10
                }
                axios({
                    url: '/job/list',
                    data: {
                        id: this.id,
                        name: this.name,
                        sort_value: "id desc",
                        offset: offset,
                        limit: 10,
                        ...params,
                    },
                    type: 'json',
                }).then(data => {
                    console.log(data)
                    const pagination = {...this.pagination};
                    // Read total count from server
                    // pagination.total = data.totalCount;
                    pagination.total = data.data.total;
                    this.loading = false;
                    this.data = data.data.data;
                    this.pagination = pagination;
                });
            },
            detail(id) {
                this.$router.push("/job/detail?id=" + id);
            },
            filter() {
                let para = [];
                if (this.id !== undefined && this.id !== '') {
                    para = [...para, 'id=' + this.id]
                }
                if (this.name !== undefined && this.name !== '') {
                    para = [...para, 'name=' + this.name]
                }
                let url = "/job/list";
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
        },
    }

</script>

<style scoped>
    .highlight {
        background-color: rgb(255, 192, 105);
        padding: 0px;
    }
</style>