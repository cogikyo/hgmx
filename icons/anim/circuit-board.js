import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const traces = [el.querySelector('path[d^="M7 21"]'), el.querySelector('path[d^="M11 9"]')];
  animate(traces, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12) });

  const pads = [el.querySelector('circle[cx="15"]'), el.querySelector('circle[cx="9"]')];
  animate(pads, { scale: [1, 1.25, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
