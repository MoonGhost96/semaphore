<template>
  <v-form
      ref="form"
      lazy-validation
      v-model="formValid"
      v-if="item != null"
  >
    <v-alert
        :value="formError"
        color="error"
        class="pb-2"
    >{{ formError }}
    </v-alert>

    <v-text-field
        v-model="item.name"
        label="秘钥名称"
        :rules="[v => !!v || '名称为必填项']"
        required
        :disabled="formSaving"
    />

    <v-select
        v-model="item.type"
        label="类型"
        :rules="[v => (!!v || !canEditSecrets) || '类型为必填项']"
        :items="inventoryTypes"
        item-value="id"
        item-text="name"
        :required="canEditSecrets"
        :disabled="formSaving || !canEditSecrets"
    />

    <v-text-field
        v-model="item.login_password.login"
        label="Login (Optional)"
        v-if="item.type === 'login_password'"
        :disabled="formSaving || !canEditSecrets"
    />

    <v-text-field
        v-model="item.login_password.password"
        label="Password"
        :rules="[v => (!!v || !canEditSecrets) || 'Password 为必填项']"
        v-if="item.type === 'login_password'"
        :required="canEditSecrets"
        :disabled="formSaving || !canEditSecrets"
        autocomplete="new-password"
    />

    <v-text-field
      v-model="item.ssh.login"
      label="Username (Optional)"
      v-if="item.type === 'ssh'"
      :disabled="formSaving || !canEditSecrets"
    />

    <!--    <v-text-field-->
    <!--        v-model="item.ssh.passphrase"-->
    <!--        label="Passphrase (Optional)"-->
    <!--        v-if="item.type === 'ssh'"-->
    <!--        :disabled="formSaving || !canEditSecrets"-->
    <!--    />-->

    <v-textarea
      outlined
      v-model="item.ssh.private_key"
      label="Private Key"
      :disabled="formSaving || !canEditSecrets"
      :rules="[v => !!v || 'Private Key 为必填项']"
      v-if="item.type === 'ssh'"
    />

    <v-checkbox
        v-model="item.override_secret"
        label="Override"
        v-if="!isNew"
    />

    <v-alert
        dense
        text
        type="info"
        v-if="item.type === 'none'"
    >
      Use this type of key for HTTPS repositories and for
      playbooks which use non-SSH connections.
    </v-alert>
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],
  data() {
    return {
      inventoryTypes: [{
        id: 'ssh',
        name: 'SSH Key',
      }, {
        id: 'login_password',
        name: 'Login with password',
      }, {
        id: 'none',
        name: 'None',
      }],
    };
  },

  computed: {
    canEditSecrets() {
      return this.isNew || this.item.override_secret;
    },
  },

  methods: {
    getNewItem() {
      return {
        ssh: {},
        login_password: {},
      };
    },

    getItemsUrl() {
      return `/api/project/${this.projectId}/keys`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/keys/${this.itemId}`;
    },
  },
};
</script>
