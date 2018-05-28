package main

import (
    "./esxi"
    "fmt"
    "errors"
    "github.com/hashicorp/terraform/helper/schema"
)

func resourceGUEST() *schema.Resource {
    return &schema.Resource{
        Create: resourceGUESTCreate,
        Read:   resourceGUESTRead,
        Update: resourceGUESTUpdate,
        Delete: resourceGUESTDelete,
        Schema: map[string]*schema.Schema{
            "esxi_hostname": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_hostname", "esxi"),
                Description: "The esxi hostname or IP address.",
            },
            "esxi_hostport": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_hostport", "22"),
                Description: "ssh port.",
            },
            "esxi_username": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_username", "root"),
                Description: "esxi ssh username.",
            },
            "esxi_password": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_password", "unset"),
                Description: "esxi ssh password.",
            },
            "clone_from_vm": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
                DefaultFunc: schema.EnvDefaultFunc("clone_from_vm", nil),
                Description: "Source vm path on esxi host to clone.",
            },
            "ovf_source": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
                DefaultFunc: schema.EnvDefaultFunc("ovf_source", nil),
                Description: "Local path to source ovf files.",
            },
            "esxi_disk_store": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_disk_store", "Least Used"),
                Description: "esxi diskstore for boot disk.",
            },
            //"esxi_virtual_network": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("esxi_virtual_network", nil),
            //    Description: "esxi virtual network.",
            //},
            "esxi_resource_pool": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
                DefaultFunc: schema.EnvDefaultFunc("esxi_resource_pool", "/"),
                Description: "Use resource pool.",
            },
            "guest_name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
                DefaultFunc: schema.EnvDefaultFunc("guest_name", "vm-example"),
                Description: "esxi guest name.",
            },
            //"guest_disk_type": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_disk_type", nil),
            //    Description: "Guest guest disk type .",
            //},
            //"guest_storage": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_storage", nil),
            //    Description: "Guest guest additional storage.",
            //},
            //"guest_nic_type": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_nic_type", nil),
            //    Description: "Guest guest nic type.",
            //},
            //"guest_mac_address": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_mac_address", nil),
            //    Description: "Guest guest mac address.",
            //},
            //"guest_memsize": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_memsize", nil),
            //    Description: "Guest guest memory size.",
            //},
            //"guest_numvcpus": &schema.Schema{
            //    Type:     schema.TypeString,
            //    Required: true,
            //    DefaultFunc: schema.EnvDefaultFunc("guest_numvcpus", nil),
            //    Description: "Guest guest number of virtual cpus.",
            //},
        },
    }
}

func resourceGUESTCreate(d *schema.ResourceData, m interface{}) error {
    esxi_hostname      := d.Get("esxi_hostname").(string)
    esxi_hostport      := d.Get("esxi_hostport").(string)
    esxi_username      := d.Get("esxi_username").(string)
    esxi_password      := d.Get("esxi_password").(string)
    clone_from_vm      := d.Get("clone_from_vm").(string)
    ovf_source         := d.Get("ovf_source").(string)
    esxi_disk_store    := d.Get("esxi_disk_store").(string)
    //esxi_virtual_network := d.Get("esxi_virtual_network").(string)
    esxi_resource_pool := d.Get("esxi_resource_pool").(string)
    guest_name         := d.Get("guest_name").(string)
    //guest_disk_type    := d.Get("guest_disk_type").(string)
    //guest_storage      := d.Get("guest_storage").(string)
    //guest_nic_type     := d.Get("guest_nic_type").(string)
    //guest_mac_address  := d.Get("guest_mac_address").(string)
    //guest_memsize      := d.Get("guest_memsize").(string)
    //guest_numvcpus     := d.Get("guest_numvcpus").(string)

    // Validations
    var src_path string
    encoded_esxi_password := esxi_password

    if clone_from_vm != "" {
      src_path = fmt.Sprintf("vi://%s:%s@%s/%s", esxi_username, encoded_esxi_password, esxi_hostname, clone_from_vm)
    } else if ovf_source != "" {
      src_path = ovf_source
    } else {
      fmt.Println("[provider-esxi] Error: You must specify clone_from_vm or src_path as a source.")
      return errors.New("Error: You must specify clone_from_vm or src_path as a source.")
    }


    vmid, err := esxi.GuestCreate(esxi_hostname, esxi_hostport, esxi_username, encoded_esxi_password,
       guest_name, esxi_disk_store, src_path, esxi_resource_pool)
    if err == 0 {
      //d.SetId(strconv.Itoa(vm_id))
      d.SetId(vmid)
    } else {
      fmt.Println("Error: Unable to create guest.")
      return errors.New(vmid)
    }
    return nil
}

func resourceGUESTRead(d *schema.ResourceData, m interface{}) error {
    return nil
}

func resourceGUESTUpdate(d *schema.ResourceData, m interface{}) error {
    return nil
}

func resourceGUESTDelete(d *schema.ResourceData, m interface{}) error {
  esxi_hostname := d.Get("esxi_hostname").(string)
  esxi_hostport := d.Get("esxi_hostport").(string)
  esxi_username := d.Get("esxi_username").(string)
  esxi_password := d.Get("esxi_password").(string)
  err := esxi.GuestDelete(esxi_hostname, esxi_hostport, esxi_username, esxi_password, d.Id())
  if err == 0 {
    d.SetId("")
  }
  return nil
}
