import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const bag = [el.querySelector('path[d^="M3.4"]'), el.querySelector('path[d^="M3.103"]')];
  animate(bag, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const handle = el.querySelector('path[d^="M16 10"]');
  animate(handle, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
