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
      label="Host Name"
      :rules="[v => !!v || 'Name 为必填项']"
      required
      :disabled="formSaving"
    />

    <v-text-field
      v-model="item.host_IP"
      label="Host IP"
      :rules="[v => !!v || 'IP 为必填项']"
      required
      :disabled="formSaving"
    />

    <v-text-field
      v-model="item.user_name"
      label="User Name"
      :rules="[v => !!v || 'User Name 为必填项']"
      required
      :disabled="formSaving"
    />

    <v-text-field
      v-model="item.password"
      label="Password"
      :rules="[v => !!v || 'Password 为必填项']"
      required
      :disabled="formSaving"
    />
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],

  methods: {
    getItemsUrl() {
      return `/api/project/${this.projectId}/host`;
    },

    getSingleItemUrl() {
      return `/api/project/${this.projectId}/host/${this.itemId}`;
    },
  },
};
</script>
