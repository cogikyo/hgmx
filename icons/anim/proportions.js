import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const splits = [el.querySelector('path[d^="M2 9h13"]'), el.querySelector('path[d="M12 9v11"]')];
  animate(splits, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
