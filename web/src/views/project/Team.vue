<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items != null">
    <EditDialog
      v-model="editDialog"
      :save-button-text="(this.itemId === 'new' ? 'Link' : 'Save')"
      :title="(this.itemId === 'new' ? 'New' : 'Edit') + ' Team Member'"
      @save="loadItems()"
    >
      <template v-slot:form="{ onSave, onError, needSave, needReset }">
        <TeamMemberForm
          :project-id="projectId"
          :item-id="itemId"
          @save="onSave"
          @error="onError"
          :need-save="needSave"
          :need-reset="needReset"
        />
      </template>
    </EditDialog>

    <YesNoDialog
      title="删除团队成员"
      text="确定删除该团队成员?"
      v-model="deleteItemDialog"
      @yes="deleteItem(itemId)"
    />

    <v-toolbar flat >
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>团队管理</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        @click="editItem('new')"
      >添加团队成员</v-btn>
    </v-toolbar>

    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
      class="mt-4"
      :items-per-page="Number.MAX_VALUE"
    >
      <template v-slot:item.admin="{ item }">

        <v-switch
          v-model="item.admin"
          inset
          :disabled="!isUserAdmin()"
          @change="item.admin ? grantAdmin(item.id) : refuseAdmin(item.id)"
        ></v-switch>
      </template>

      <template v-slot:item.actions="{ item }">
        <v-btn
          icon
          :disabled="!isUserAdmin()"
          @click="askDeleteItem(item.id)"
        >
          <v-icon>mdi-delete</v-icon>
        </v-btn>
      </template>
    </v-data-table>
  </div>

</template>
<script>
import ItemListPageBase from '@/components/ItemListPageBase';
import TeamMemberForm from '@/components/TeamMemberForm.vue';
import axios from 'axios';

export default {
  components: { TeamMemberForm },
  mixins: [ItemListPageBase],
  methods: {
    async grantAdmin(userId) {
      await axios({
        method: 'post',
        url: `/api/project/${this.projectId}/users/${userId}/admin`,
        responseType: 'json',
      });
      await this.loadItems();
    },
    async refuseAdmin(userId) {
      await axios({
        method: 'delete',
        url: `/api/project/${this.projectId}/users/${userId}/admin`,
        responseType: 'json',
      });
      await this.loadItems();
    },
    getHeaders() {
      return [
        {
          text: '姓名',
          value: 'name',
          width: '50%',
        },
        {
          text: '系统用户名',
          value: 'username',
        },
        {
          text: '邮箱',
          value: 'email',
          width: '50%',
        },
        {
          text: '管理员',
          value: 'admin',
        },
        {
          text: '操作',
          value: 'actions',
          sortable: false,
        }];
    },
    getSingleItemUrl() {
      return `/api/project/${this.projectId}/users/${this.itemId}`;
    },
    getItemsUrl() {
      return `/api/project/${this.projectId}/users?sort=name&order=asc`;
    },
    getEventName() {
      return 'i-repositories';
    },
    isUserAdmin() {
      return ((this.items.find((x) => x.id === this.userId) || {}).admin || localStorage.getItem('userAdmin') === 'true');
    },
  },
};
</script>
