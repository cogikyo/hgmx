import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const slice = el.querySelector('path[d^="m7.2"]');
  animate(slice, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const cherry = el.querySelector("circle");
  animate(cherry, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
