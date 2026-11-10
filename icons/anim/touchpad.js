import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pad = el.querySelector("rect");
  animate(pad, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const divider = el.querySelector('path[d="M2 14h20"]');
  animate(divider, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });

  const split = el.querySelector('path[d="M12 20v-6"]');
  animate(split, { scaleY: [0, 1] }, { duration: τmin, delay: τmid });
}
