import { animate, τmax } from "../motion.js";

export default function (el) {
  const scribble = el.querySelector('path[d^="m21 17"]');
  animate(scribble, { scaleX: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
