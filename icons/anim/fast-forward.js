import { animate, stagger } from "../motion.js";

export default function (el) {
  const heads = el.querySelectorAll("path");
  animate(heads, { translateX: [0, 3, -1, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
