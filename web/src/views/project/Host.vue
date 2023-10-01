<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items != null">
    <EditDialog
      v-model="editDialog"
      :save-button-text="itemId === 'new' ? '创建' : '保存'"
      :title="`${itemId === 'new' ? '创建' : '编辑'}主机`"
      :max-width="450"
      position="top"
      @save="loadItems()"
    >
      <template v-slot:form="{ onSave, onError, needSave, needReset }">
        <HostForm
          :project-id="projectId"
          :item-id="itemId"
          @save="onSave"
          @error="onError"
          :need-save="needSave"
          :need-reset="needReset"
        />
      </template>
    </EditDialog>

    <ObjectRefsDialog
      object-title="主机"
      :object-refs="itemRefs"
      :project-id="projectId"
      v-model="itemRefsDialog"
    />

    <YesNoDialog
      title="删除主机"
      text="确认删除该主机?"
      v-model="deleteItemDialog"
      @yes="deleteItem(itemId)"
    />

    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>主机管理</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        @click="editItem('new')"
      >创建主机
      </v-btn>
    </v-toolbar>

    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
      class="mt-4"
      :items-per-page="Number.MAX_VALUE"
    >
      <template v-slot:item.type="{ item }">
        <code>{{ item.type }}</code>
      </template>
      <template v-slot:item.actions="{ item }">
        <div style="white-space: nowrap">
          <v-btn
            icon
            class="mr-1"
            @click="askDeleteItem(item.id)"
          >
            <v-icon>mdi-delete</v-icon>
          </v-btn>

          <v-btn
            icon
            class="mr-1"
            @click="editItem(item.id)"
          >
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
        </div>
      </template>
    </v-data-table>
  </div>

</template>
<script>
import ItemListPageBase from '@/components/ItemListPageBase';
import HostForm from '@/components/HostForm.vue';

export default {
  components: { HostForm },
  mixins: [ItemListPageBase],
  methods: {
    getHeaders() {
      return [{
        text: 'Name',
        value: 'name',
        width: '25%',
      },
      {
        text: 'IP',
        value: 'host_ip',
        width: '25%',
      },
      {
        text: 'UserName',
        value: 'user_name',
        width: '25%',
      },
      {
        text: 'Password',
        value: 'password',
        width: '25%',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      }];
    },
    getItemsUrl() {
      return `/api/project/${this.projectId}/host`;
    },
    getSingleItemUrl() {
      return `/api/project/${this.projectId}/host/${this.itemId}`;
    },
    getEventName() {
      return 'i-host';
    },
  },
};
</script>
