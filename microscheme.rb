require 'bigdecimal'
require 'ostruct'

class LexBinding
  attr_accessor :parent, :local

  def initialize(parent = nil)
    @parent = parent
    @local = {}
  end

  def set(k, v)
    local[k] = v
  end

  def get(k)
    return local[k] if local.key?(k)
    return parent.get(k) if parent
    nil
  end

  def dump
    current = self
    indent = ''
    res = ''

    until current.nil?
      current.local.each do |k, v|
        res += indent + "#{k}: #{v}\n"
      end
      indent += '  '
      current = current.parent
    end

    res
  end
end

def tokenize(src)
  tokens = []
  current_val = ''

  src.each_char do |c|
    case c
    when '(', ')'
      unless current_val.empty?
        tokens.push current_val
        current_val = ''
      end
      tokens.push c
    when ' ', "\n", "\t", "\r"
      unless current_val.empty?
        tokens.push current_val
        current_val = ''
      end
    else
      current_val += c
    end
  end

  tokens.push(current_val) unless current_val.empty?
  tokens
end

def parse(tokens)
  compound_stack = []
  root = nil

  tokens.each do |t|
    case t
    when '('
      expr = OpenStruct.new(children: [], compound: true)

      if root.nil?
        root = expr
      elsif !compound_stack.empty?
        compound_stack.last.children.push(expr)
      else
        raise "Invalid token in sequence: #{t}"
      end

      compound_stack.push(expr)
    when ')'
      raise "Invalid ')'" if compound_stack.empty?
      compound_stack.pop
    else
      expr = OpenStruct.new(token: t)

      if root.nil?
        root = expr
      elsif !compound_stack.empty?
        compound_stack.last.children.push(expr)
      else
        raise "Invalid token in sequence: #{t}"
      end
    end
  end

  raise "Unclosed compound expression" unless compound_stack.empty?

  root
end

SPECIAL_FORMS = ['let', 'if', 'lambda', '+', '=', '-']

def classify(expr)
  if expr.compound
    first = expr.children.first
    if first.nil?
      :null
    elsif first.compound
      :application
    elsif SPECIAL_FORMS.include?(first.token)
      first.token.to_sym
    else
      :application
    end
  elsif expr.token =~ /^\d+$/
    :number
  elsif expr.token == '#t'
    :true
  elsif expr.token == '#f'
    :false
  else
    :reference
  end
end

def evaluate(expr, lexbinding)
  case classify(expr)
  when :null
    OpenStruct.new(type: :null)
  when :number
    v = BigDecimal.new(expr.token)
    if v.frac == 0
      v = v.to_i
    else
      v = v.to_f
    end

    OpenStruct.new(type: :number, underlying: v)
  when :+
    underlying = expr.children[1..-1].reduce(0) do |memo, c|
      val = evaluate(c, lexbinding)
      raise "Invalid operand for + expression: #{val}" if val.type != :number
      memo + val.underlying
    end

    OpenStruct.new(type: :number, underlying: underlying)
  when :-
    if expr.children.size != 3
      raise "Invalid subtraction: #{expr}"
    end

    vals = expr.children[1..-1].map { |c| evaluate(c, lexbinding) }
    if vals.any? { |v| v.type != :number }
      raise "Invalid operand for - expression: #{val}"
    end

    OpenStruct.new(type: :number, underlying: vals.first.underlying - vals.last.underlying)
  when :'='
    vals = expr.children[1..-1].map { |c| evaluate(c, lexbinding) }
    first = vals.first
    vals[1..-1].each do |c|
      if first.type != c.type || first.underlying != c.underlying
        return OpenStruct.new(type: :boolean, underlying: false)
      end
    end

    OpenStruct.new(type: :boolean, underlying: true)
  when :true
    OpenStruct.new(type: :boolean, underlying: true)
  when :false
    OpenStruct.new(type: :boolean, underlying: false)
  when :reference
    v = lexbinding.get(expr.token)
    raise "Invalid reference: #{expr}" unless v
    v
  when :let
    if expr.children.size != 3 || !expr.children[1].compound
      raise "Invalid let expression: #{expr}"
    end

    next_binding = LexBinding.new(lexbinding)
    expr.children[1].children.each do |c|
      if !c.compound || c.children.size != 2 || c.children[0].compound
        raise "Invalid let expression assignment: #{c}"
      end

      next_binding.set(
        c.children[0].token,
        evaluate(c.children[1], lexbinding)
      )
    end

    evaluate(expr.children[2], next_binding)
  when :lambda
    if expr.children.size != 3 || !expr.children[1].compound || !expr.children[2].compound
      raise "Invalid lambda expression: #{expr}"
    end

    if expr.children[1].children.any? { |c| c.compound }
      raise "Invalid lambda parameter list: #{expr.children[1]}"
    end

    OpenStruct.new(
      type: :function,
      params: expr.children[1].children.map { |c| c.token },
      body: expr.children[2],
      lexbinding: lexbinding,
    )
  when :application
    f = evaluate(expr.children[0], lexbinding)
    if f.type != :function
      raise "Invalid function: #{expr.children[0]}"
    end

    if f.params.size != expr.children.size - 1
      raise "Invalid function invocation: #{expr.children.size - 1} of #{f.params.size} arguments provided"
    end

    next_binding = LexBinding.new(f.lexbinding)
    expr.children[1..-1].each_with_index do |c, i|
      next_binding.set(f.params[i], evaluate(c, lexbinding))
    end

    evaluate(f.body, next_binding)
  when :if
    if expr.children.size != 4
      raise "Invalid if expression: #{expr}"
    end

    pred = evaluate(expr.children[1], lexbinding)
    if pred.type != :boolean
      raise "Predicate #{expr.children[1]} evaluates to invalid value: #{pred}"
    end

    if pred.underlying == true
      evaluate(expr.children[2], lexbinding)
    else
      evaluate(expr.children[3], lexbinding)
    end
  end
end

def interpret(src)
  toks = tokenize(src)

  begin
    syntree = parse(toks)
  rescue => e
    puts e
    puts toks
    return
  end

  begin
    evaluate(syntree, LexBinding.new)
  rescue => e
    puts e
    puts syntree
  end
end

# puts interpret('1')
# puts interpret('(- 100 1)')
# puts interpret('(- 0 1)')
# puts interpret('(+ 1 (+ 2 3))')

# puts interpret('(= 1 1)')
# puts interpret('(= 1 2)')
# puts interpret('(= 1 1 1)')
# puts interpret('(= 1 1 2)')
# puts interpret('(= #t #t)')
# puts interpret('(= #f #f)')
# puts interpret('(= #t #f)')

# puts interpret('(let ((x (+ 1 5)) (y 4)) (+ y x))')

# puts interpret('((lambda (x y) (+ x y)) 1 2)')

# puts interpret(%q|
# (let ((add1 (lambda (x) (+ x 1)))
#       (add2 (lambda (x) (+ x 2)))
#       (x 5))
#   (add1 (add2 (add1 x))))
# |)

# puts interpret('(if #t 4 5)')
# puts interpret('(if #f 4 5)')
# puts interpret('(if (= 1 1) 4 5)')
# puts interpret('(if (= 1 2) 4 5)')
# puts interpret('(if (= 1 2) (+1 2) (+ 3 4))')

puts interpret(%q|
(let ((tri (lambda (x f)
             (if (= x 0)
                 0
                 (+ x (f (- x 1) f))))))
  (tri 100 tri))
|)
