<template>
  <v-form
    ref="form"
    lazy-validation
    v-model="formValid"
  >
    <v-alert
      :value="formError"
      color="error"
      class="pb-2"
    >{{ formError }}</v-alert>

    <v-text-field
      v-model="item.password"
      label="Password"
      :type="showPassword ? 'text' : 'password'"
      :append-icon="showPassword ? 'mdi-eye' : 'mdi-eye-off'"
      @click:append="showPassword = !showPassword"
      :rules="[v => !!v || 'Password 为必填项']"
      required
      :disabled="formSaving"
    ></v-text-field>
  </v-form>
</template>
<script>
import ItemFormBase from '@/components/ItemFormBase';

export default {
  mixins: [ItemFormBase],

  data() {
    return {
      showPassword: false,
    };
  },

  methods: {
    async loadData() {
      this.item = {};
    },

    getItemsUrl() {
      return null;
    },

    getSingleItemUrl() {
      return null;
    },

    getRequestOptions() {
      return {
        method: 'post',
        url: `/api/users/${this.itemId}/password`,
      };
    },
  },
};
</script>
