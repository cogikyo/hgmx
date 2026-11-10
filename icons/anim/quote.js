import { animate, stagger } from "../motion.js";

export default function (el) {
  const marks = [el.querySelector('path[d^="M5 3"]'), el.querySelector('path[d^="M16 3"]')];
  animate(marks, { translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
