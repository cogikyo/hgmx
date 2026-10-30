import { animate, stagger } from "../motion.js";

export default function (el) {
  const chevrons = [el.querySelector('path[d^="m17 11"]'), el.querySelector('path[d^="m17 18"]')];
  animate(chevrons, { translateY: [0, -3, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
