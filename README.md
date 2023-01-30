# WebResume

<img align="right" width="159px" style="margin: 20px;" src="https://raw.githubusercontent.com/sculley/web-resume/main/web-resume-logo.png">

WebResume is a Go-based application designed to act as a dynamic home page and resume for engineers. It utilizes templating to populate values in the resume and retrieves these values from a YAML file, making it easy for users to change and customize their information. With WebResume, engineers can showcase their skills, experience, and projects cleanly and professionally, making it an essential tool for job seekers also it provides a clean homepage for your domain. The application's user-friendly interface, combined with its flexible and customizable design, makes it an ideal solution for anyone looking to make a strong first impression in their job search.

## Getting Started

The easiest way to get started with WebResume is to use the Docker image. This image is available on Docker Hub and can be pulled with the following command:

```bash
docker pull yelluc/webresume
```

Once the image is pulled, you can run it with the following command:

```bash
docker run -d -p 8080:8080 -v /path/to/your/config:/usr/local/web-resume/config --name web-resume sculley/webresume
```

This will run the container in the background, exposing port 8080 on the host machine. It will also mount the config directory from the host machine to the container, allowing you to edit the config file without having to rebuild the image. The container will be named `web-resume` and can be stopped and started with the following commands:

```bash
docker stop web-resume
docker start web-resume
```

You can also run the application directly from the source code. To do this, you will need to have Go installed on your machine. Once Go is installed, you can run the following commands to get the application running:

```bash
git clone git@github.com:sculley/web-resume.git
cd web-resume
go run cmd/web-resume/main.go
```

## Usage

Once the application is running, you can access it by navigating to `http://localhost:8080` in your browser. The application will display a home page with your name and a subtitle. Clicking on the About link in the navigation bar will take you to your resume page, which will display your name, subtitle, experience, skills, and certifications.

You can change the port that the application runs on by setting the `PORT` environment variable. For example, if you wanted to run the application on port 80, you would run the following command:

```bash
PORT=80 go run cmd/web-resume/main.go
```

Below are the environment variables that can be set, along with their default values:

| Environment Variable | Value |
| --- | --- |
| CONFIG_PATH | /usr/local/web-resume/config |
| GIN_MODE | release |
| LOG_FORMAT | json |
| LOG_LEVEL | info |
| PORT | 8080 |
| STATIC_PATH | /usr/local/web-resume/static |
| TEMPLATES_PATH | /usr/local/web-resume/templates |

## Configuration

The configuration file is a YAML file that contains all of the information that will be displayed on the resume. The following is an example of the configuration file:

```yaml
title: "Sam Culley"
full_name: "Sam Culley"
github_profile: "https://github.com/sculley"
home_page_subtitle: "I'm a Software Engineer"
about_page_subtitle: "My name is Sam Culley and I'm a Software Engineer"
experience:
  - company: "Acme Inc."
    url: "https://acme.com"
    role: "Software Engineer"
    start_date: "January 2020"
    end_date: "January 2022"
    accomplishments:
      - "Built a new feature for the product"
      - "Fixed a bug in the product"
skills:
  - Go
certifications:
  - Example Certification
```

## License

This project is licensed under the MIT License - see the [LICENSE](MIT-LICENSE.txt) file for details
