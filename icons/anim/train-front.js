import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const train = el.querySelector("g");
  animate(train, { scale: [0.5, 1.1, 1], opacity: [0, 1] }, { duration: τmax });

  const rails = el.querySelectorAll('path[d^="m8 19"], path[d^="m16 19"]');
  animate(rails, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const lights = el.querySelectorAll('path[d^="m9 15"], path[d^="m15 15"]');
  animate(lights, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
