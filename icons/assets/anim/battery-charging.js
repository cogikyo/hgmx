import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = [el.querySelector('path[d^="M14.856"]'), el.querySelector('path[d^="M5.14"]')];
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const bolt = el.querySelector('path[d^="m11 7"]');
  animate(bolt, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const tip = el.querySelector('path[d="M22 14v-4"]');
  animate(tip, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
