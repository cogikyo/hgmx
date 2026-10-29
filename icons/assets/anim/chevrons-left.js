import { animate, stagger } from "../motion.js";

export default function (el) {
  const chevrons = [el.querySelector('path[d^="m11 17"]'), el.querySelector('path[d^="m18 17"]')];
  animate(chevrons, { translateX: [0, -3, 1, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
