import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticks = el.querySelectorAll("path");
  animate(ticks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });
  animate(ticks, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
