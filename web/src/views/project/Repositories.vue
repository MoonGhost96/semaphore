<template xmlns:v-slot="http://www.w3.org/1999/XSL/Transform">
  <div v-if="items != null && keys != null">
    <EditDialog
      v-model="editDialog"
      :save-button-text="itemId === 'new' ? '创建' : '保存'"
      :title="`${itemId === 'new' ? '创建' : '编辑'}仓库`"
      @save="loadItems()"
      max-width="500"
    >
      <template v-slot:form="{ onSave, onError, needSave, needReset }">
        <RepositoryForm
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
      object-title="仓库"
      :object-refs="itemRefs"
      :project-id="projectId"
      v-model="itemRefsDialog"
    />

    <YesNoDialog
      title="删除仓库"
      text="确定删除?"
      v-model="deleteItemDialog"
      @yes="deleteItem(itemId)"
    />

    <v-toolbar flat >
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>仓库管理</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn
        color="primary"
        @click="editItem('new')"
      >创建仓库</v-btn>
    </v-toolbar>

    <v-data-table
      :headers="headers"
      :items="items"
      hide-default-footer
      class="mt-4"
      :items-per-page="Number.MAX_VALUE"
    >
      <template v-slot:item.git_url="{ item }">
        {{ item.git_url }}
        <code v-if="!item.git_url.startsWith('/')">{{ item.git_branch }}</code>
      </template>

      <template v-slot:item.ssh_key_id="{ item }">
        {{ keys.find((k) => k.id === item.ssh_key_id).name }}
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
import RepositoryForm from '@/components/RepositoryForm.vue';
import axios from 'axios';

export default {
  mixins: [ItemListPageBase],
  components: { RepositoryForm },
  data() {
    return {
      keys: null,
    };
  },

  async created() {
    this.keys = (await axios({
      method: 'get',
      url: `/api/project/${this.projectId}/keys`,
      responseType: 'json',
    })).data;
  },

  methods: {
    getHeaders() {
      return [{
        text: '名称',
        value: 'name',
        width: '25%',
      },
      {
        text: 'Git URL',
        value: 'git_url',
        width: '50%',
      },
      {
        text: 'SSH 秘钥',
        value: 'ssh_key_id',
        width: '25%',
      },
      {
        text: '操作',
        value: 'actions',
        sortable: false,
      }];
    },
    getItemsUrl() {
      return `/api/project/${this.projectId}/repositories`;
    },
    getSingleItemUrl() {
      return `/api/project/${this.projectId}/repositories/${this.itemId}`;
    },
    getEventName() {
      return 'i-repositories';
    },
  },
};
</script>
