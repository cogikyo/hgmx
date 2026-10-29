import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M10 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const shackle = el.querySelector('path[d^="M20 17"]');
  animate(shackle, { translateY: [-3, 1, 0] }, { duration: τmid, delay: τmin });

  const body = el.querySelector("rect");
  animate(body, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
