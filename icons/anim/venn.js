import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const sets = el.querySelectorAll("circle");
  animate(sets, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid, delay: stagger(0.12) });
}
