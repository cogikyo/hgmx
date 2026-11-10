import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const nut = el.querySelector("g");
  animate(nut, { translateY: [-2, 0.5, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d="m2 2 20 20"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
