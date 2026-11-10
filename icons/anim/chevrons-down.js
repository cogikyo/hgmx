import { animate, stagger } from "../motion.js";

export default function (el) {
  const chevrons = [el.querySelector('path[d^="m7 13"]'), el.querySelector('path[d^="m7 6"]')];
  animate(chevrons, { translateY: [0, 3, -1, 0] }, { duration: 0.3, delay: stagger(0.12) });
}
