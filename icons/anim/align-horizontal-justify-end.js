import { animate, stagger } from "../motion.js";

export default function (el) {
  const row = [el.querySelector('rect[x="12"]'), el.querySelector('rect[x="2"]')];
  animate(row, { translateX: [0, 1.5, -0.5, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
