import { animate, stagger } from "../motion.js";

export default function (el) {
  const grips = [el.querySelector('path[d^="m21 15"]'), el.querySelector('path[d^="M21 8"]')];
  animate(grips, { translateX: [0, 2, -0.5, 0], translateY: [0, 2, -0.5, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
