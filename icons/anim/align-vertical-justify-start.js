import { animate, stagger } from "../motion.js";

export default function (el) {
  const row = [el.querySelector('rect[y="6"]'), el.querySelector('rect[y="16"]')];
  animate(row, { translateY: [0, -1.5, 0.5, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
