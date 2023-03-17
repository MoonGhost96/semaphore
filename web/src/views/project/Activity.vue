<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items">
    <v-toolbar flat >
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>仪表盘</v-toolbar-title>
      <v-spacer></v-spacer>
      <div>
        <v-tabs centered>
          <v-tab key="history" :to="`/project/${projectId}/history`">运行历史</v-tab>
          <v-tab key="activity" :to="`/project/${projectId}/activity`">操作记录</v-tab>
          <v-tab key="settings" :to="`/project/${projectId}/settings`">项目设置</v-tab>
        </v-tabs>
      </div>
    </v-toolbar>

    <v-data-table
      :headers="headers"
      :items="items"
      class="mt-4"
      :footer-props="{ itemsPerPageOptions: [20] }"
    >
      <template v-slot:item.created="{ item }">
        {{ item.created | formatDate }}
      </template>
    </v-data-table>
  </div>
</template>
<script>
import ItemListPageBase from '@/components/ItemListPageBase';

export default {
  mixins: [ItemListPageBase],

  methods: {
    getHeaders() {
      return [
        {
          text: '时间',
          value: 'created',
          sortable: false,
          width: '20%',
        },
        {
          text: '用户',
          value: 'username',
          sortable: false,
          width: '10%',
        },
        {
          text: '详情',
          value: 'description',
          sortable: false,
          width: '70%',
        },
      ];
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/events/last`;
    },
  },
};
</script>
