<template>
  <div id="app">
    <v-app id="inspire">
      <v-data-table
        :headers="headers"
        :items="campaigns"
        sort-by="calories"
        class="elevation-1"
      >
        <template v-slot:item.Status="{ item }">
          <v-chip :color="getColor(item.Status)" dark>
            {{ item.status }}
          </v-chip>
        </template>
        <template v-slot:top>
          <v-toolbar flat color="white">
            <v-toolbar-title>My Campaigns </v-toolbar-title>
            <v-divider class="mx-4" inset vertical></v-divider>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialog" max-width="500px">
              <template v-slot:activator="{ on }">
                <v-btn color="primary" dark class="mb-2" v-on="on"
                  >New Item</v-btn
                >
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">{{ formTitle }}</span>
                </v-card-title>

                <v-card-text>
                  <v-container>
                    <v-row>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                          v-model="editedItem.name"
                          label="Name"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                          v-model="editedItem.status"
                          label="Status"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                          v-model="editedItem.type"
                          label="Type"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" sm="6" md="4">
                        <v-text-field
                          v-model="editedItem.budget"
                          label="Budget"
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </v-container>
                </v-card-text>

                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" text @click="close"
                    >Cancel</v-btn
                  >
                  <v-btn color="blue darken-1" text @click="save">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-toolbar>
        </template>
        <template v-slot:item.action="{ item }">
          <v-icon small class="mr-2" @click="editItem(item)"> edit </v-icon>
          <v-icon small @click="deleteItem(item)"> delete </v-icon>
        </template>
        <template v-slot:no-data>
          <v-btn color="primary" @click="initialize">Reset</v-btn>
        </template>
      </v-data-table>
    </v-app>
  </div>
</template>

<script>
import axios from "axios";
export default {
  data: () => ({
    dialog: false,
    headers: [
      {
        text: "Name",
        align: "left",
        sortable: false,
        value: "name",
      },
      { text: "Status", value: "status" },
      { text: "Type", value: "type" },
      { text: "Budget", value: "budget" },
      { text: "Created", value: "created_on" },
      { text: "Actions", value: "action", sortable: false },
    ],
    campaigns: [],
    editedIndex: -1,
    editedItem: {
      name: "",
      status: "",
      type: "",
      budget: 0,
    },
    defaultItem: {
      name: "",
      status: "",
      type: "",
      budget: 0,
    },
  }),

  computed: {
    formTitle() {
      return this.editedIndex === -1 ? "New Item" : "Edit Item";
    },
  },

  watch: {
    dialog(val) {
      val || this.close();
    },
  },

  created() {
    this.initialize();
  },

  methods: {
    async initialize() {
      let { data } = await axios.get("http://localhost:8000/api/campaign");
      this.campaigns = data;
    },

    editItem(item) {
      this.editedIndex = this.campaigns.indexOf(item);
      this.editedItem = Object.assign({}, item);
      this.dialog = true;
    },

    deleteItem(item) {
      const index = this.campaigns.indexOf(item);
      confirm("Are you sure you want to delete this item?") &&
        this.campaigns.splice(index, 1);
    },

    close() {
      this.dialog = false;
      setTimeout(() => {
        this.editedItem = Object.assign({}, this.defaultItem);
        this.editedIndex = -1;
      }, 300);
    },

    save() {
      if (this.editedIndex > -1) {
        Object.assign(this.campaigns[this.editedIndex], this.editedItem);
      } else {
        this.editedItem.budget = parseInt(this.editedItem.budget);
        this.campaigns.push(this.editedItem);

        axios.post(
          "http://localhost:8000/api/newcampaign",
          JSON.stringify(this.editedItem)
        );
      }
      this.close();
    },

    getColor(calories) {
      if (calories) return "red";
      return "orange";
    },
  },
};
</script>
