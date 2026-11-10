import { animate, stagger } from "../motion.js";

export default function (el) {
  const bands = el.querySelectorAll('rect:not([x="10"]), path');
  animate(bands, { scaleX: [0, 1.1, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
