# Ixia-C Operations

The easiest way to manage emulated network topologies involving [Ixia-C](https://github.com/open-traffic-generator/ixia-c).

### Getting Started

- Install latest

    ```sh
    go install github.com/open-traffic-generator/ixops@latest
    ```

- Check default configuration

    ```sh
    ixops config get
    ```

- Setup emulated topology as specified in configuration

    ```sh
    ixops topology create
    ```

- Generate test UDP traffic

    ```sh
    ixops otg --pps 100 --count 500 --udp gen
    ```

- Teardown topology as specified in configuration

    ```sh
    ixops topology delete
    ```

- Check usage

    ```sh
    ixops help
    ```
