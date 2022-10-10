from ast import main
import datetime
import os
import psycopg2
import redis

#configurations and processes
try:
    # create users_files directory if it does not exist
    # this directory stores individual users files



    #connect to the database
    r = redis.Redis(decode_responses=True)
    connection = psycopg2.connect(user="jugo",
                                  password="jugo",
                                  host="127.0.0.1",
                                  port="5432",
                                  database="jugo_db")

    # get cursor to for query
    cursor = connection.cursor()
    print("connected")
    #this directory contains the temporay files to be appended to the user files
    temporary_user_files_directory = "/home/bethel/Base/projects/jugo/jugo_api/uploads"
    print("temp")
    # loop to continuesly check redis list(tasks queue),and compare it to the 
    # files in the temporary storage,if it matches ,then read its content, and 
    # add it to the users file

    while True:
      # get the items in the task queue, and get the number of elements in it
      task_list = r.lrange("tasks_queue",0,-1)
      items_in_queue = len(task_list)

      # if there is an item or task in the queue,pop the last on the list(this is a blocking operations)
      # the task item is a tuple that contains the name of the list and the name of the file to be processed
      #  get the name of the the file to be processed ,then iterate over every file in the temporary directory
      #  comparing the filename the file,if any of the file matches ,read it ,and add the contents to the 
      # appropriate user(this can be gotten from the filename>>"username_filename.txt"),and delete temporary file
      if items_in_queue >= 1 :
        print(f"this is the tasks list::{task_list}")
        print(f"there are {items_in_queue} tasks in queue")

        task = r.brpop("tasks_queue",timeout=0)
        file_from_task = task[1]
        print(f"task>>{task}")
        #iterate over every file in the temporary directory
        for filename in os.listdir(temporary_user_files_directory):
          if file_from_task == filename:
            # the files come in the form 'username_filename.txt',split the filename at '_'
            # to a list containing the username and filename
            seperated = filename.split("_")
            print(f"seperated user and file>>{seperated}")
            #username from splitting the filename
            user = seperated[0]

            #if the present index of the for loop is a file,read its contents
            user_temporary_file = os.path.join(temporary_user_files_directory,filename)
            if os.path.isfile(user_temporary_file):
              try:
                user_file = open(user_temporary_file,"r")
                contents = user_file.read()
              except:
                print("error could not open or read file")

            #query the database,for the user whose name is gotten from the filename splitted
            cursor.execute("SELECT name,file_path,last_changed,email FROM users_with_filespath WHERE name=%s",(user,))
  
            user_query = cursor.fetchone()
            print(user_query)
            #get the file path from the query,if its ""(empty string),create a file using the user's name
            #and append the content of the temporary file(file from temporary directory with users name)
            username=user_query[0]
            user_file_path=user_query[1]
            main_directory = os.getcwd()
            if not os.path.exists("../users_files"):
              os.mkdir("../users_files")
            print(f"main_directory::{main_directory}")
            if user_file_path == "":
              user_permanent_file = username + ".txt"
              os.chdir("../")
              pwd = os.getcwd()
              print(pwd)
              path = os.path.join(pwd,"users_files",user_permanent_file)

              with open(path,"x") as f:
                try:
                  f.write("\n"+contents)
                  print("creating user file")
                  print("file contents appended to user file")
                except:
                  print("Could not open or append to user permanent file")

              update_query = """Update users_with_filespath set file_path=%s,last_changed=%s where name=%s"""
              cursor.execute(update_query,(path,datetime.datetime.now(),username))
              connection.commit()
              os.chdir(main_directory)
              present_directory = os.getcwd()
              print(f"present_directory::{present_directory}")

            #if filepath exists,apend contents of temporary file to the file
            else:
              with open(user_file_path,"a") as f:
                f.write("\n"+contents)
                print(user_file_path)
                print("file contents appended to user file")

            #delete temporary file
            os.remove(user_temporary_file)
            print("temporary file deleted")
            print("------------------------------------------")

    
#error while connecting to the database
except (Exception, psycopg2.Error) as error:
    print("Error while connecting to PostgreSQL", error)
#close connection when done         
finally:
    if connection:
        cursor.close()
        connection.close()
        print("PostgreSQL connection is closed")
