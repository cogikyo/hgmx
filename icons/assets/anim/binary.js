import { animate, stagger } from "../motion.js";

export default function (el) {
  const [one, other] = el.querySelectorAll("g");
  const digits = [el.querySelector('rect[x="6"]'), one, other, el.querySelector('rect[x="14"]')];
  animate(digits, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
