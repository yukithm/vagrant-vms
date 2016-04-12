#!/usr/bin/env ruby
# encoding: utf-8

require "open3"

class VagrantBox

end

class VagrantVM
  GLOBAL_STATUS_REGEXP = /\A(?<id>[0-9a-f]+)\s+(?<name>\S+)\s+(?<provider>\S+)\s+(?<state>\S+)\s+(?<directory>.+?)\s*\z/
  attr_reader :id, :name, :provider, :state, :directory,
              :vm_id, :index_uuid, :exist
  alias_method :exist?, :exist

  def initialize(id:, name:, provider:, state:, directory:)
    @id = id
    @name = name
    @provider = provider
    @state = state
    @directory = directory
    if Dir.exist?(machine_dir)
      @vm_id = read_vm_id
      @index_uuid = read_index_uuid
      @exist = true
    else
      @exist = false
    end
  end

  def self.global_status
    o, s = Open3.capture2("vagrant", "global-status")
    fail "Failed to execute vagrant global-status" unless s.success?

    vms = []
    o.each_line do |line|
      line.chomp!
      m = GLOBAL_STATUS_REGEXP.match(line)
      next unless m

      vms << new(
        id: m[:id],
        name: m[:name],
        provider: m[:provider],
        state: m[:state],
        directory: m[:directory],
      )
    end
    vms
  end

  def machine_dir
    File.join(directory, ".vagrant", "machines", name, provider)
  end

  private

  def read_index_uuid
    File.read(File.join(machine_dir, "index_uuid"))
  end

  def read_vm_id
    File.read(File.join(machine_dir, "id"))
  end
end

vms = VagrantVM.global_status
vms.each do |vm|
  # p vm
  puts "#{vm.directory} (#{vm.vm_id})"
end
