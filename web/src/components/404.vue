<template>
    <div>
        <a-row type="flex" justify="space-around" align="middle">
            <a-col><img src="../assets/404.png" :style="{display: 'inline-block', height: 'auto', 'max-width': '100%'}">
            </a-col>
            <a-progress type="circle" status="active" :percent="percent"
                        :style="{ position: 'fixed', zIndex: 1, width: '100%' }">
                <template v-slot:format="percent">
                    <span style="color: red">{{percent}}</span>
                </template>
            </a-progress>
        </a-row>

    </div>
</template>
<script>
    export default {
        name: "System404",

        data() {
            return {
                percent: 0
            };
        },

        mounted() {
            this.setPercent();
        },

        methods: {
            setPercent() {
                if (this.percent < 100) {
                    this.percent = this.percent + 1;
                } else if (this.percent === 100 || this.percent > 100) {
                    this.$router.push({
                        path: "/"
                    });
                }
            },

            timer() {
                return setTimeout(() => {
                    this.setPercent();
                }, 50)
            },
        },

        watch: {
            percent() {
                this.timer()
            }
        },

        destroyed() {
            clearTimeout(this.timer())
        },
    };
</script>
