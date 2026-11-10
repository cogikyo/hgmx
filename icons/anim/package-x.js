import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const box = ["M12 22V12", "M21 1", "M3.29 7", "m7.5 4.27"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const bars = [el.querySelector('path[d="m16.5 14.5 5 5"]'), el.querySelector('path[d="m16.5 19.5 5-5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmin, delay: stagger(τmin, { startDelay: τmin }) });
}
