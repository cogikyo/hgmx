import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = [el.querySelector('path[d="m6 4 6 8 6-8"]'), el.querySelector('path[d="M12 12v8"]')];
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="M8 13h8"]'), el.querySelector('path[d="M8 17h8"]')];
  animate(bars, { scaleX: [0, 1.25, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
